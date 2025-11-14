package initialize

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
)

// ES 初始化Elasticsearch
func ES() error {
	client, err := initElasticsearchClient()
	if err != nil {
		return err
	}
	global.GVA_ES = client

	// 如果配置了默认索引，创建它
	if global.GVA_CONFIG.Elasticsearch.Index != "" {
		if err := CreateDefaultIndex(client, global.GVA_CONFIG.Elasticsearch.Index); err != nil {
			global.GVA_LOG.Warn("创建默认索引失败", zap.Error(err))
		}
	}

	return nil
}

// initElasticsearchClient 初始化Elasticsearch客户端
func initElasticsearchClient() (*elasticsearch.Client, error) {
	cfg := global.GVA_CONFIG.Elasticsearch

	if len(cfg.Addresses) == 0 {
		return nil, fmt.Errorf("elasticsearch addresses cannot be empty")
	}

	// 创建Elasticsearch配置
	esConfig := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	}

	// 配置重试
	if cfg.MaxRetries > 0 {
		esConfig.MaxRetries = cfg.MaxRetries
	}

	// 配置超时时间
	if cfg.Timeout > 0 {
		timeout := time.Duration(cfg.Timeout) * time.Second
		esConfig.Transport = &http.Transport{
			Dial: (&net.Dialer{
				Timeout: timeout,
			}).Dial,
			MaxIdleConnsPerHost:   100,
			ResponseHeaderTimeout: timeout,
		}
	}

	// 配置HTTPS和证书
	if len(cfg.CertFile) > 0 {
		certs, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.CertFile)
		if err != nil {
			global.GVA_LOG.Error("加载HTTPS证书失败", zap.Error(err))
			return nil, err
		}

		esConfig.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
				Certificates:       []tls.Certificate{certs},
			},
		}
	}

	// 配置gzip压缩
	if cfg.Compression == "gzip" {
		esConfig.CompressRequestBody = true
	}

	// 创建Elasticsearch客户端
	client, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		global.GVA_LOG.Error("创建Elasticsearch客户端失败", zap.Error(err))
		return nil, err
	}

	// 测试连接
	res, err := client.Info()
	if err != nil {
		global.GVA_LOG.Error("Elasticsearch连接失败", zap.Error(err))
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		global.GVA_LOG.Error(fmt.Sprintf("Elasticsearch响应错误: %s", res.String()))
		return nil, fmt.Errorf("elasticsearch error: %d", res.StatusCode)
	}

	global.GVA_LOG.Info("Elasticsearch连接成功",
		zap.Strings("addresses", cfg.Addresses),
		zap.String("index", cfg.Index),
	)

	return client, nil
}

// CreateDefaultIndex 创建默认索引
func CreateDefaultIndex(client *elasticsearch.Client, indexName string) error {
	// 检查索引是否存在
	res, err := client.Indices.Exists([]string{indexName})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		global.GVA_LOG.Info(fmt.Sprintf("索引 %s 已存在", indexName))
		return nil
	}

	// 创建索引映射
	mapping := `{
        "settings": {
            "number_of_shards": 1,
            "number_of_replicas": 0,
            "analysis": {
                "analyzer": {
                    "default": {
                        "type": "standard"
                    }
                }
            }
        },
        "mappings": {
            "properties": {
                "@timestamp": {
                    "type": "date"
                },
                "message": {
                    "type": "text"
                },
                "level": {
                    "type": "keyword"
                },
                "logger": {
                    "type": "keyword"
                }
            }
        }
    }`

	// 创建索引
	res, err = client.Indices.Create(
		indexName,
		client.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		global.GVA_LOG.Error("创建Elasticsearch索引失败", zap.Error(err))
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch create index error: %s", res.String())
	}

	global.GVA_LOG.Info(fmt.Sprintf("Elasticsearch索引 %s 创建成功", indexName))
	return nil
}

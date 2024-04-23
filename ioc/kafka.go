package ioc

import (
	"github.com/IBM/sarama"
	"github.com/MuxiKeStack/be-search/events"
	"github.com/MuxiKeStack/be-search/pkg/saramax"
	"github.com/spf13/viper"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	saramaCfg.Producer.Partitioner = sarama.NewConsistentCRCHashPartitioner
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	client, err := sarama.NewClient(cfg.Addrs, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitConsumers(MySQLBinlog *events.MySQLBinlogConsumer) []saramax.Consumer {
	return []saramax.Consumer{
		MySQLBinlog,
	}
}

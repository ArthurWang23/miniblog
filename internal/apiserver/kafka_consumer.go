package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/ArthurWang23/miniblog/internal/pkg/log"
	"github.com/segmentio/kafka-go"
)

// StartKafkaConsumer 会根据 cfg.KafkaOptions 启动一个最小示例消费者协程。
// 返回的 stop 函数用于优雅停止消费者：取消上下文并关闭 reader。
// 如未配置 Kafka（broker 列表为空或 topic 为空），将返回一个空操作的 stop 函数。
func StartKafkaConsumer(cfg *Config) (func(), error) {
	// 可选：未配置直接返回 no-op
	if cfg == nil || cfg.KafkaOptions == nil || len(cfg.KafkaOptions.Brokers) == 0 || cfg.KafkaOptions.Topic == "" {
		return func() {}, nil
	}

	// 创建 Dialer（含 TLS / SASL）
	dialer, err := cfg.KafkaOptions.Dialer()
	if err != nil {
		return func() {}, err
	}

	ro := cfg.KafkaOptions.ReaderOptions
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:       cfg.KafkaOptions.Brokers,
		Topic:         cfg.KafkaOptions.Topic,
		GroupID:       ro.GroupID,
		Partition:     ro.Partition, // 注意：GroupID 与 Partition 不能同时设置，KafkaOptions.Validate 已校验
		QueueCapacity: ro.QueueCapacity,
		MinBytes:      ro.MinBytes,
		MaxBytes:      ro.MaxBytes,
		MaxWait:       ro.MaxWait,
		Dialer:        dialer,
	})
	ctx, cancel := context.WithCancel(context.Background())

	// 消费协程
	go func() {
		log.Infow("Kafka consumer started",
			"brokers", cfg.KafkaOptions.Brokers,
			"topic", cfg.KafkaOptions.Topic,
			"groupID", ro.GroupID,
			"partition", ro.Partition,
		)
		defer log.Infow("Kafka consumer stopped")

		for {
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				// 若短暂的网络或 broker 问题，打印错误并继续
				log.Errorw("Kafka consumer read error", "err", err)
				time.Sleep(300 * time.Millisecond)
				continue
			}

			// 尝试按 JSON 解码，否则打印原始值
			var obj map[string]any
			if err := json.Unmarshal(msg.Value, &obj); err == nil {
				log.Infow("Kafka message",
					"topic", msg.Topic,
					"partition", msg.Partition,
					"offset", msg.Offset,
					"headers", msg.Headers,
					"value.json", obj,
				)
			} else {
				log.Infow("Kafka message",
					"topic", msg.Topic,
					"partition", msg.Partition,
					"offset", msg.Offset,
					"headers", msg.Headers,
					"value.raw", string(msg.Value),
				)
			}
		}
	}()

	// 返回停止函数
	return func() {
		cancel()
		_ = reader.Close()
	}, nil
}
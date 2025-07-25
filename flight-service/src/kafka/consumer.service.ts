import { Injectable, OnApplicationShutdown } from '@nestjs/common';
import {
  Consumer,
  ConsumerRunConfig,
  ConsumerSubscribeTopics,
  Kafka,
} from 'kafkajs';

@Injectable()
export class ConsumerService implements OnApplicationShutdown {
  constructor() {}

  private readonly kafka = new Kafka({
    clientId: 'flight-service',
    brokers: ['kafka:9092'],
  });
  private readonly consumers: Consumer[] = [];

  async consumeMessage(
    topic: ConsumerSubscribeTopics,
    config: ConsumerRunConfig,
  ) {
    const consumer = this.kafka.consumer({ groupId: 'booking-group' });
    await consumer.connect();
    await consumer.subscribe(topic);
    await consumer.run(config);
    this.consumers.push(consumer);
  }

  async onApplicationShutdown(signal: string) {
    console.log(`Received shutdown signal: ${signal}`);
    for (const consumer of this.consumers) {
      await consumer.disconnect();
    }
  }
}

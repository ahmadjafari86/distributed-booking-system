import {
  Injectable,
  OnApplicationShutdown,
  OnModuleInit,
} from '@nestjs/common';
import { Kafka, Producer } from 'kafkajs';

@Injectable()
export class ProducerService implements OnModuleInit, OnApplicationShutdown {
  constructor() {}
  private readonly kafka = new Kafka({
    clientId: 'flight-service',
    brokers: ['kafka:9092'],
  });

  private readonly producer: Producer = this.kafka.producer();

  async onModuleInit() {
    await this.producer.connect();
  }

  async sendMessage(topic: string, message: any) {
    try {
      await this.producer.send({
        topic,
        messages: [
          {
            value: JSON.stringify(message),
          },
        ],
      });
      console.log(`Message sent to topic ${topic}:`, message);
    } catch (error) {
      console.error(`Failed to send message to topic ${topic}:`, error);
    }
  }

  async onApplicationShutdown(signal: string) {
    console.log(`Received shutdown signal: ${signal}`);
    await this.producer.disconnect();
  }
}

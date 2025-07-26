import {
  Injectable,
  InternalServerErrorException,
  OnApplicationShutdown,
  OnModuleInit,
} from '@nestjs/common';
import { Kafka, Producer } from 'kafkajs';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class ProducerService implements OnModuleInit, OnApplicationShutdown {
  private readonly kafka: Kafka;
  private readonly producer: Producer;
  private readonly KAFKA_BROKERS: string[];
  private readonly KAFKA_CLIENT_ID: string;

  constructor(private readonly configService: ConfigService) {
    const brokersString = this.configService.get<string>('KAFKA_BROKERS');
    if (!brokersString) {
      throw new InternalServerErrorException(
        'KAFKA_BROKERS environment variable is not set.',
      );
    }
    this.KAFKA_BROKERS = brokersString.split(',');

    const clientId = this.configService.get<string>('KAFKA_CLIENT_ID');
    if (!clientId) {
      throw new InternalServerErrorException(
        'KAFKA_CLIENT_ID environment variable is not set.',
      );
    }
    this.KAFKA_CLIENT_ID = clientId;

    this.kafka = new Kafka({
      clientId: this.KAFKA_CLIENT_ID,
      brokers: this.KAFKA_BROKERS,
    });

    this.producer = this.kafka.producer();
  }

  async onModuleInit() {
    console.log('ProducerService: Connecting Kafka producer...');
    await this.producer.connect();
    console.log('ProducerService: Kafka producer connected.');
  }

  async sendMessage(topic: string, key: string, message: any) {
    try {
      await this.producer.send({
        topic,
        messages: [
          {
            key: key,
            value: JSON.stringify(message),
          },
        ],
      });
      console.log(`ProducerService: Message sent to topic ${topic}:`, message);
    } catch (error) {
      console.error(
        `ProducerService: Failed to send message to topic ${topic}:`,
        error,
      );
    }
  }

  async onApplicationShutdown(signal: string) {
    console.log(
      `ProducerService: Received shutdown signal: ${signal}. Disconnecting Kafka producer.`,
    );
    await this.producer.disconnect();
    console.log('ProducerService: Kafka producer disconnected.');
  }
}

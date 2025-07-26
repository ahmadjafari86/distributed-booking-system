import {
  BadRequestException,
  Injectable,
  InternalServerErrorException,
  NotFoundException,
} from '@nestjs/common';
import { CreateFlightDto } from './dto/create-flight.dto';
import { UpdateFlightDto } from './dto/update-flight.dto';
import { ReserveSeatDto } from './dto/reserve-seat.dto';
import { InjectRepository } from '@nestjs/typeorm';
import { Flight } from './entities/flight.entity';
import { Repository } from 'typeorm';
import { FlightReservation } from './entities/flight-reservation.entity';
import { ConfigService } from '@nestjs/config';
import { ProducerService } from 'src/kafka/producer.service';

@Injectable()
export class FlightsService {
  KAFKA_REQUEST_TOPIC: string;
  KAFKA_RESULT_TOPIC: string;

  constructor(
    @InjectRepository(Flight)
    private readonly flightRepository: Repository<Flight>,
    @InjectRepository(FlightReservation)
    private reservationRepository: Repository<FlightReservation>,
    private readonly configService: ConfigService,
    private readonly producerService: ProducerService,
  ) {
    const kafkaRequestTopic = this.configService.get<string>(
      'KAFKA_REQUEST_TOPIC',
    );
    const kafkaResultTopic =
      this.configService.get<string>('KAFKA_RESULT_TOPIC');
    if (!kafkaRequestTopic || !kafkaResultTopic) {
      throw new InternalServerErrorException(
        'KAFKA_BROKERS or KAFKA_CLIENT_ID environment variable is not set.',
      );
    }
    this.KAFKA_REQUEST_TOPIC = kafkaRequestTopic;
    this.KAFKA_RESULT_TOPIC = kafkaResultTopic;
  }

  async create(createFlightDto: CreateFlightDto) {
    const flight = this.flightRepository.create(createFlightDto);
    return await this.flightRepository.save(flight);
  }

  async findAll() {
    const flights = await this.flightRepository.find({
      relations: ['reservations'],
    });

    return flights.map((flight) => {
      const reservedSeats = flight.reservations.reduce(
        (total, reservation) => total + reservation.seatCount,
        0,
      );

      return {
        id: flight.id,
        flightNumber: flight.flightNumber,
        departure: flight.departure,
        arrival: flight.arrival,
        totalSeats: flight.totalSeats,
        availableSeats: flight.totalSeats - reservedSeats,
        reservations: flight.reservations.map((res) => ({
          id: res.id,
          seatCount: res.seatCount,
        })),
      };
    });
  }

  async findOne(id: string) {
    return await this.flightRepository.findOne({ where: { id } });
  }

  async update(id: string, updateFlightDto: UpdateFlightDto) {
    const flight = await this.flightRepository.findOne({ where: { id } });
    if (!flight) {
      throw new NotFoundException(`Flight with ID ${id} not found`);
    }
    Object.assign(flight, updateFlightDto);
    return await this.flightRepository.save(flight);
  }

  async remove(id: string) {
    const flight = await this.flightRepository.findOne({ where: { id } });
    if (!flight) {
      throw new NotFoundException(`Flight with ID ${id} not found`);
    }
    return await this.flightRepository.remove(flight);
  }

  async reserveSeats(reserveDto: ReserveSeatDto) {
    const flight = await this.flightRepository.findOne({
      where: { id: reserveDto.flightId },
      relations: ['reservations'],
    });

    if (!flight) {
      throw new NotFoundException(
        `Flight with ID ${reserveDto.flightId} not found`,
      );
    }

    if (reserveDto.seatCount <= 0) {
      throw new BadRequestException('Seat count must be a positive integer');
    }

    const reservedSeats = flight.reservations.reduce(
      (sum, r) => sum + r.seatCount,
      0,
    );
    const availableSeats = flight.totalSeats - reservedSeats;

    if (reserveDto.seatCount > availableSeats) {
      throw new BadRequestException('Not enough available seats');
    }

    const reservation = this.reservationRepository.create({
      flight,
      seatCount: reserveDto.seatCount,
    });

    const savedReservation = await this.reservationRepository.save(reservation);

    await this.producerService.sendMessage(
      this.KAFKA_RESULT_TOPIC,
      reservation.id,
      {
        flightId: reserveDto.flightId,
        seatCount: reserveDto.seatCount,
        reservationId: reservation.id,
      },
    );

    return savedReservation;
  }

  async cancelReservation(id: string) {
    const reservation = await this.reservationRepository.findOne({
      where: { id },
      relations: ['flight'],
    });

    if (!reservation) {
      throw new NotFoundException('Reservation not found');
    }

    await this.reservationRepository.remove(reservation);
    return 'Reservation canceled successfully';
  }

  async getFlightReservations(flightId: string) {
    const flight = await this.findOne(flightId);
    if (!flight) {
      throw new NotFoundException(`Flight with ID ${flightId} not found`);
    }

    return this.reservationRepository.find({
      where: { flight: { id: flight.id } },
    });
  }
}

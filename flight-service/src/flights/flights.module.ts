import { Module } from '@nestjs/common';
import { FlightsService } from './flights.service';
import { FlightsController } from './flights.controller';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Flight } from './entities/flight.entity';
import { FlightReservation } from './entities/flight-reservation.entity';

@Module({
  imports: [TypeOrmModule.forFeature([Flight, FlightReservation])],
  controllers: [FlightsController],
  providers: [FlightsService],
})
export class FlightsModule {}

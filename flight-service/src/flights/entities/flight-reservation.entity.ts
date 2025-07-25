import { Entity, PrimaryGeneratedColumn, Column, ManyToOne } from 'typeorm';
import { Flight } from './flight.entity';

@Entity()
export class FlightReservation {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  seatCount: number;

  @ManyToOne(() => Flight, (flight) => flight.reservations, {
    onDelete: 'CASCADE',
  })
  flight: Flight;
}

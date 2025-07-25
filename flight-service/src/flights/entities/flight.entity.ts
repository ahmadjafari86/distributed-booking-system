import { Column, Entity, OneToMany, PrimaryGeneratedColumn } from 'typeorm';
import { FlightReservation } from './flight-reservation.entity';

@Entity({ name: 'flights' })
export class Flight {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ type: 'varchar', length: 100, unique: true })
  flightNumber: string;

  @Column({ type: 'varchar', length: 100 })
  departure: string;

  @Column({ type: 'varchar', length: 100 })
  arrival: string;

  @Column({ type: 'integer', nullable: false })
  totalSeats: number;

  @OneToMany(() => FlightReservation, (res) => res.flight)
  reservations: FlightReservation[];
}

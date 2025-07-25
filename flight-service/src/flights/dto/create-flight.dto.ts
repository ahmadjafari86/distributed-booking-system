import { IsPositive, IsString } from 'class-validator';

export class CreateFlightDto {
  @IsPositive()
  @IsString()
  flightNumber: string;

  @IsString()
  departure: string;

  @IsString()
  arrival: string;

  @IsPositive()
  total_seats: number;
}

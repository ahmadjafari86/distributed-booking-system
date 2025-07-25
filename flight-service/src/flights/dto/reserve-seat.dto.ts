import { IsUUID, IsInt } from 'class-validator';

export class ReserveSeatDto {
  @IsUUID()
  flightId: string;

  @IsInt()
  seatCount: number;
}

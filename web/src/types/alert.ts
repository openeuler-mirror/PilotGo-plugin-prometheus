export interface Alert {
  id: number;
  ip: string;
  alertName: string;
  level: string;
  alertTime: string;
  alertEndTime: string;
  handleState: string;
  summary: string;
  description: string;
}
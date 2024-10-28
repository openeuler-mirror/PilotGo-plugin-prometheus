export interface ConfigRule {
  id?: number;
  metrics: string;
  ips?: string[];
  batches: any[];
  batches_str: string;
  alertTargets: any[];
  alertHostIds?: string;
  alertLabel: string;
  alertName: string;
  duration: number | string;
  severity: string;
  input_severity?: string;
  threshold: number | string;
  desc: string;
  [key: string]: unknown;
}

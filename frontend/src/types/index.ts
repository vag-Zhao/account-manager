export interface Account {
  id: number;
  account: string;
  password: string;
  accountType: 'PLUS' | 'BUSINESS' | 'FREE';
  isSold: boolean;
  soldAt: string | null;
  expireAt: string | null;
  reminderSent: boolean;
  notes: string;
  createdAt: string;
  updatedAt: string;
}

export interface AccountFilter {
  accountType: string;
  isSold: boolean | null;
  search: string;
  page: number;
  pageSize: number;
}

export interface AccountStats {
  total: number;
  plusCount: number;
  businessCount: number;
  freeCount: number;
  soldCount: number;
  expiredCount: number;
  expiringIn7Days: number;
}

export interface PaginatedAccounts {
  data: Account[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

export interface EmailConfig {
  id: number;
  smtpHost: string;
  smtpPort: number;
  senderEmail: string;
  senderPassword: string;
  recipientEmail: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface EmailLog {
  id: number;
  subject: string;
  content: string;
  recipient: string;
  status: string;
  error: string;
  createdAt: string;
}

export interface SystemConfig {
  id: number;
  defaultValidityDays: number;
  reminderDaysBefore: number;
  copyFormat: string;
  emailFormat: string;
  accountTypes: string;
  accountStatuses: string;
  createdAt: string;
  updatedAt: string;
}

export interface Tag {
  label: string;
  value: string;
  color: string;
}

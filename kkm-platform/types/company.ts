/**
 * Company types and interfaces for the companies management system
 */

export interface Company {
  id: string;
  name: string;
  description: string;
  owner_id: string;
  member_count: number;
  status: string;
  created_at: number; // Unix timestamp
  updated_at: number; // Unix timestamp
}

export interface CompanyResponse {
  company?: Company;
  companies?: Company[];
  data?: Company | Company[];
  message?: string;
  error?: string;
}

export interface CreateCompanyRequest {
  name: string;
  description: string;
}

export interface UpdateCompanyRequest {
  name?: string;
  description?: string;
  status?: string;
}

export interface User {
  id?: string | number;
  username: string;
  email: string;
  password?: string;
  phone?: string;
  nickname?: string;
  role_id: number;
  role?: {
    id: number;
    name: string;
  };
  status: number;
  class_id?: number;
  created_at?: string;
  updated_at?: string;
}

export interface User {
  id?: string;
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
  created_at?: string;
  updated_at?: string;
}

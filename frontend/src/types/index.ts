export interface Book {
  id: string;
  title: string;
  owner_id: string;
  cover_photo_id: string | null;
  bg_color: string;
  created_at: string;
  updated_at: string;
}

export interface BookDetail extends Book {
  photos?: Photo[];
}

export interface Photo {
  id: string;
  book_id: string;
  filename: string;
  caption: string | null;
  position: number;
  width: number;
  height: number;
  taken_at: string | null;
}

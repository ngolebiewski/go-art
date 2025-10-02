// src/data/mockArtworks.ts
import { type Artwork } from '../interfaces/artwork';

export const ARTWORK_DATA: Artwork[] = [
  {
    id: 1,
    title: "Birch Bark",
    thumbnailUrl: "/images/demo-art/birch.jpg", // Use a path relative to your public folder
    ownerName: "Nick G."
  },
  {
    id: 2,
    title: "Lichen branch from the Adirondacks",
    thumbnailUrl: "/images/demo-art/lichen_branch.jpg",
    ownerName: "Nick G."
  },
  {
    id: 3,
    title: "Lichen and Branch Segment",
    thumbnailUrl: "/images/demo-art/lichen_flame.jpg",
    ownerName: "Nick G."
  },
];
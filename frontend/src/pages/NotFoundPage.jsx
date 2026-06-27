import { Link } from "react-router-dom";
import { Home } from "lucide-react";

export default function NotFoundPage() {
  return (
    <div className="min-h-[70vh] flex items-center justify-center px-4">
      <div className="text-center">
        <p className="text-8xl font-black text-gray-200 mb-4">404</p>
        <h1 className="text-2xl font-bold text-gray-900 mb-2">
          Halaman Tidak Ditemukan
        </h1>
        <p className="text-gray-500 mb-6">
          Halaman yang Anda cari tidak tersedia atau telah dipindahkan.
        </p>
        <Link to="/" className="btn-primary inline-flex items-center gap-2">
          <Home className="w-4 h-4" /> Kembali ke Beranda
        </Link>
      </div>
    </div>
  );
}

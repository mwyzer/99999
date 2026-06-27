export default function Footer() {
  return (
    <footer className="bg-gray-800 text-gray-300 mt-auto">
      <div className="max-w-7xl mx-auto px-4 py-8 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div>
            <h3 className="text-white font-semibold text-lg mb-3">
              PropertyHub
            </h3>
            <p className="text-sm">
              Platform multi-tenant untuk listing properti dari agensi, bank,
              dan perusahaan terpercaya di Indonesia.
            </p>
          </div>
          <div>
            <h4 className="text-white font-medium mb-3">Tautan</h4>
            <ul className="space-y-1 text-sm">
              <li>
                <a href="/" className="hover:text-white transition-colors">
                  Cari Properti
                </a>
              </li>
              <li>
                <a
                  href="/register"
                  className="hover:text-white transition-colors"
                >
                  Daftar
                </a>
              </li>
              <li>
                <a href="/login" className="hover:text-white transition-colors">
                  Masuk
                </a>
              </li>
            </ul>
          </div>
          <div>
            <h4 className="text-white font-medium mb-3">Kontak</h4>
            <p className="text-sm">Email: support@propertyhub.id</p>
            <p className="text-sm mt-2">
              © {new Date().getFullYear()} PropertyHub. All rights reserved.
            </p>
          </div>
        </div>
      </div>
    </footer>
  );
}

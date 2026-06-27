import { Loader2 } from "lucide-react";

export default function Loading({ fullScreen }) {
  const content = (
    <div className="flex items-center justify-center gap-3 text-gray-500">
      <Loader2 className="w-6 h-6 animate-spin text-primary-600" />
      <span>Memuat...</span>
    </div>
  );

  if (fullScreen) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        {content}
      </div>
    );
  }

  return (
    <div className="py-20 flex items-center justify-center">{content}</div>
  );
}

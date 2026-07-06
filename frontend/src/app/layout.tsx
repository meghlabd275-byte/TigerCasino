import type { Metadata } from 'next';
import './globals.css';
import ClientProviders from './ClientProviders';
import { Toaster } from 'react-hot-toast';

export const metadata: Metadata = {
  title: 'TigerCasino - Crypto Casino Platform',
  description: 'The ultimate cryptocurrency casino experience with instant transactions and provably fair gaming',
  keywords: ['crypto casino', 'bitcoin casino', 'online gambling', 'provably fair'],
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className="bg-tiger-dark text-white font-body antialiased">
        <ClientProviders>
          {children}
          <Toaster 
            position="top-right"
            toastOptions={{
              style: {
                background: '#16213E',
                color: '#fff',
                border: '1px solid #FF6B35',
              },
            }}
          />
        </ClientProviders>
      </body>
    </html>
  );
}

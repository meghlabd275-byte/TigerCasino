import type { Metadata } from 'next';
import '@/styles/globals.css';
import ClientProviders from './ClientProviders';

export const metadata: Metadata = {
  title: 'TigerCasino - Crypto Casino Platform',
  description: 'The ultimate cryptocurrency casino experience with instant transactions and fair play',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body>
        <ClientProviders>
          {children}
        </ClientProviders>
      </body>
    </html>
  );
}

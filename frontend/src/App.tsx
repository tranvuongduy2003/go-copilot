import { ErrorBoundary } from '@/components/common/error-boundary';
import { AppProviders } from '@/providers';
import { router } from '@/router';
import { RouterProvider } from 'react-router-dom';

function App() {
  return (
    <ErrorBoundary>
      <AppProviders>
        <RouterProvider router={router} />
      </AppProviders>
    </ErrorBoundary>
  );
}

export default App;

import AppLayout from 'components/layouts/AppLayout';
import { NextPageWithLayout } from 'pages/_app';


const App: NextPageWithLayout = () => {
  return (
    <main>
      Welcome to Beans.
    </main>
  )
}

App.getLayout = page => <AppLayout>{page}</AppLayout>

export default App

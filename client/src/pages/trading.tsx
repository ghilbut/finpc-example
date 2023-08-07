import { trpc } from '~/utils/trpc'

export default function Home() {
  const { data } = trpc.getStockList.useQuery();

  return (
    <main>
        <h1>Test</h1>
        <table>
            <thead>
                <tr>
                    <td>ID</td>
                    <td>Code</td>
                    <td>Name</td>
                    <td>Total Stock Count</td>
                </tr>
            </thead>
            <tbody>
            {data && data.map(stock => {
                return (
                    <tr key="{stock.id}">
                        <td>{stock.id}</td>
                        <td>{stock.code}</td>
                        <td>{stock.name}</td>
                        <td>{stock.totalStockCount.toLocaleString()}</td>
                    </tr>
                );
            })}
            </tbody>
        </table>
    </main>
  )
}

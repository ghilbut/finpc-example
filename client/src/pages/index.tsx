import Link from 'next/link';
import { trpc } from '~/utils/trpc'

export default function Home() {
  const { data } = trpc.listSubjects.useQuery();

  return (
    <main>
        <h1>Subject List</h1>
        <ul>
            {data && data.map(subject => {
                return (<li key="{subject.id}"><Link href={`/${subject.id}`}>{subject.title}</Link></li>);
            })}
        </ul>
    </main>
  )
}

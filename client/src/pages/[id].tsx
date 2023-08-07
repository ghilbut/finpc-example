import Link from "next/link";
import { useRouter } from 'next/router'
import { trpc } from '~/utils/trpc'

export default function Home() {
    const router = useRouter();

    const id = (router.query.id || '0').toString()
    const subject = trpc.getSubject.useQuery({ id: parseInt(id)});
    const list = trpc.listQuestions.useQuery();

    return (
        <main>
            <h1>Question List</h1>
            <h2 key={subject.data && subject.data.id}>{subject.data && subject.data.title}</h2>
            <ul>
                {list.data && list.data.map(q => {
                    return (<li key="{q.id}">{q.question}</li>);
                })}
            </ul>
        </main>
    )
}

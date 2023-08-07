import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router'
import { trpc } from '~/utils/trpc'

export default function Home() {
    const router = useRouter();
    const id = parseInt((router.query.id || '0').toString());

    const [ newQuestion, setNewQuestion ] = useState('');

    function onNewQuestion(q) {

    }

    const subject = trpc.getSubject.useQuery({ id: id });
    const list = trpc.listQuestions.useQuery({ id: id });

    const like = trpc.like.useMutation();
    const unlike = trpc.like.useMutation();
    const createQuestion = trpc.createQuestion.useMutation();
    const deleteQuestion = trpc.deleteQuestion.useMutation();

    async function doRefresh() {
        console.log('doRefresh');
        await list.refetch();
    }

    async function doLike(id: number) {
        console.log(`doLike(qid: ${id})`);
        await like.mutate({ id });
        await list.refetch();
    }

    async function doCreateQuestion() {
        console.log(`doCreateQuestion(q: ${newQuestion})`);
        await like.mutate({ id });
        await list.refetch();
    }

    return (
        <main>
            <h1>Question List</h1>
            <h2 key={subject.data && subject.data.id}>{subject.data && subject.data.title}</h2>
            <div>
                <button onClick={doRefresh}>Refresh</button>
                <ul>
                    {list.data && list.data.map(q => {
                        return (<li key="{q.id}">
                            {q.question}&nbsp;&nbsp;
                            ({q.like})&nbsp;&nbsp;
                            <button onClick={ ()=>{doLike(q.id)} }>Like</button>
                        </li>);
                    })}
                </ul>
            </div>
            <div>
                <input type="text" onChange={ (e)=>{setNewQuestion(e.target.value)} } value={newQuestion}/>
                <button onClick={doCreateQuestion}>New Question</button>
            </div>
        </main>
    )
}

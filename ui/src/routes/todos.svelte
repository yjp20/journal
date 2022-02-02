<script context="module">
	import { api } from '$lib/api'

	export async function load({ fetch, session }) {
		return {
			props: {
				todos: await api('GET', fetch, session, 'todo')
			}
		}
	}
</script>

<script>
	import { flip } from 'svelte/animate'
	import Todo from '$lib/Todo.svelte'
	import TodoEdit from '$lib/TodoEdit.svelte'

	export let todos
	let edit = undefined

	$: sorted = todos.sort((a, b) => {
		if (a.completed !== b.completed) return a.completed - b.completed
		if (a.completed_date !== b.completed_date)
			return new Date(b.completed_date) - new Date(a.completed_date)
		if (a.cart !== b.cart) return b.cart - a.cart
		if (a.due_date !== b.due_date) {
			const a_date = a.due_date !== null ? new Date(a.due_date) : Infinity
			const b_date = b.due_date !== null ? new Date(b.due_date) : Infinity
			return a_date - b_date
		}
		return a.id - b.id
	})
</script>

<svelte:head>
	<title>journal | todos</title>
</svelte:head>

<h1 class="title">todos</h1>
<div class="todonew">
	<TodoEdit bind:todos />
</div>
{#if sorted.length > 0 && !sorted[0].cart}
	<p class="paragraph"><em>No todos in cart</em></p>
{/if}
{#each sorted as todo (todo.id)}
	<div
		animate:flip={{ duration: 100 }}
		class:is-completed={todo.completed}
		class:is-cart={todo.cart}
		class="todo-animate todo-wrapper"
	>
		<Todo bind:todos bind:todo bind:edit />
	</div>
{/each}

<style>
	.todo-wrapper + .todo-wrapper {
		margin-top: 0.25rem;
	}

	.todo-wrapper.is-cart:not(.is-completed) + .todo-wrapper:not(.is-cart):not(.is-completed) {
		margin-top: 1rem;
	}

	.todo-wrapper:not(.is-completed) + .todo-wrapper.is-completed {
		margin-top: 1rem;
	}

	.todonew {
		margin-bottom: 1.5em;
	}
</style>

<script>
	import { addDays } from 'date-fns'

	import { api, requireAuth } from '$lib/api'
	import { session } from '$app/stores'

	import Tag from '$lib/Tag.svelte'
	import IconButton from '$lib/IconButton.svelte'
	import TodoEdit from '$lib/TodoEdit.svelte'

	import EditIcon from '$lib/icons/edit.svelte'
	import CartIcon from '$lib/icons/shopping-cart.svelte'
	import DeleteIcon from '$lib/icons/trash.svelte'
	import ChecklistItem from '$lib/ChecklistItem.svelte'

	export let todo
	export let todos
	export let edit = undefined

	function startEdit() {
		if (requireAuth($session)) return;
		edit = todo.id
	}

	async function todoDelete() {
		if (requireAuth($session)) return;
		await api('DELETE', fetch, $session, `todo/${todo.id}`)
		todos = todos.filter((v) => v.id != todo.id)
	}

	async function todoToggleComplete(e) {
		if (requireAuth($session)) return;
		if (todo.recur) {
			todos = [
				...todos,
				await api('POST', fetch, $session, 'todo', {
					description: todo.description,
					due_date: addDays(new Date(todo.due_date), todo.recur),
					recur: todo.recur,
					private: todo.private
				})
			]
		}

		todo.completed = e.target.checked || e.currentTarget.checked
		todo.cart = false
		todo.recur = undefined
		todo.completed_date = todo.completed ? new Date() : null
		await api('PUT', fetch, $session, `todo/${todo.id}`, todo)
	}

	async function todoCart() {
		if (requireAuth($session)) return;
		todo.cart = todo.cart ? false : true
		await api('PUT', fetch, $session, `todo/${todo.id}`, todo)
	}

	function dueColor(date) {
		if (new Date(date) < new Date()) return { fg: 'var(--white)', bg: 'var(--red)' }
		if (new Date(date) < addDays(new Date(), 1)) return { fg: 'var(--white)', bg: 'var(--orange)' }
		if (new Date(date) > addDays(new Date(), 7)) return { fg: 'var(--grey)', bg: 'var(--light)' }
	}
</script>

{#if edit === todo.id}
	<div class="todoedit box">
		<TodoEdit {todo} on:finish={() => (edit = false)} />
	</div>
{:else}
	<ChecklistItem checked={todo.completed} on:change={todoToggleComplete}>
		{todo.description}

		<slot slot="tags">
			{#if todo.private}
				<Tag value={"private"} />
			{/if}
			{#if todo.blocked}
				<Tag value={"blocked"} />
			{/if}
			{#if !todo.completed && todo.due_date}
				<Tag value={new Date(todo.due_date)} {...dueColor(todo.due_date)} />
			{/if}
			{#if todo.recur}
				<Tag value="↻ {todo.recur % 7 == 0 ? `${todo.recur / 7} week` : `${todo.recur} days`}" />
			{/if}
			{#if todo.completed && todo.completed_date}
				<Tag value={new Date(todo.completed_date)} />
			{/if}
		</slot>

		<slot slot="actions">
			<IconButton on:click={startEdit} description="Edit todo"><EditIcon /></IconButton>
			<IconButton on:click={() => todoDelete(todo.id)} description="Delete todo">
				<DeleteIcon />
			</IconButton>
			{#if !todo.completed}
				<IconButton on:click={todoCart} description="Cart todo"><CartIcon /></IconButton>
			{/if}
		</slot>
	</ChecklistItem>
{/if}

<style>
	.todoedit {
		width: 100%;
		box-sizing: border-box;
		margin: 0.5em 0;
	}
</style>

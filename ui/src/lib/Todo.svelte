<script>
	import { api } from '$lib/api';
	import { session } from '$app/stores';

	import Tag from '$lib/Tag.svelte';
	import IconButton from '$lib/IconButton.svelte';
	import TodoEdit from '$lib/TodoEdit.svelte';
	import Checkbox from '$lib/Checkbox.svelte';

	import EditIcon from '$lib/icons/edit.svelte';
	import CartIcon from '$lib/icons/shopping-cart.svelte';
	import DeleteIcon from '$lib/icons/trash.svelte';

	export let todo;
	export let todos;
	export let edit = undefined;

	function startEdit() {
		edit = todo.id;
	}

	async function todoDelete() {
		await api('DELETE', fetch, $session, `todo/${todo.id}`);
		todos = todos.filter((v) => v.id != todo.id);
	}

	async function todoToggleComplete(e) {
		todo.completed = e.currentTarget.checked;
		todo.cart = false;
		todo.completed_date = todo.completed ? new Date() : null;
		await api('PUT', fetch, $session, `todo/${todo.id}`, todo);
	}

	async function todoCart() {
		todo.cart = todo.cart ? false : true;
		await api('PUT', fetch, $session, `todo/${todo.id}`, todo);
	}
</script>

<div class="todo" class:is-editing={edit === todo.id} class:is-completed={todo.completed}>
	{#if edit === todo.id}
		<div class="todoedit box">
			<TodoEdit {todo} bind:edit />
		</div>
	{:else}
		<div class="todo-toggle">
			<Checkbox
				description="Mark todo as {todo.completed ? 'not done' : 'done'}"
				checked={todo.completed}
				on:change={todoToggleComplete}
			/>
		</div>
		<div class="todo-content">
			<div class="todo-display">{todo.description}</div>
			<div class="todo-tags">
				{#if !todo.completed && todo.due_date}
					<Tag
						value={new Date(todo.due_date)}
						fg={new Date(todo.due_date) < new Date() ? 'var(--red)' : 'black'}
					/>
				{/if}
				{#if todo.completed && todo.completed_date}
					<Tag value={new Date(todo.completed_date)} />
				{/if}
			</div>
		</div>
		<div class="todo-actions">
			<IconButton on:click={startEdit} description="Edit todo"><EditIcon /></IconButton>
			<IconButton on:click={() => todoDelete(todo.id)} description="Delete todo">
				<DeleteIcon />
			</IconButton>
			{#if !todo.completed}
				<IconButton on:click={todoCart} description="Cart todo"><CartIcon /></IconButton>
			{/if}
		</div>
	{/if}
</div>

<style>
	.todo {
		display: flex;
		width: 100%;
	}

	.todo.is-completed {
		opacity: 0.5;
	}

	.todo-toggle {
		margin-top: 0.1em;
	}

	.todo-content {
		display: block;
		margin-left: 0.25em;
	}

	.todo-display {
		display: inline;
	}

	.todo-tags {
		margin-left: 0.25rem;
	}

	.todo-actions {
		display: flex;
		margin-left: auto;
		opacity: 0;
		align-self: flex-start;
	}

	.todo:hover .todo-actions {
		opacity: 1;
	}

	.todo:hover .todo-display {
		text-decoration: underline;
	}

	.todoedit {
		width: 100%;
		margin: 0.5em 0;
	}
</style>

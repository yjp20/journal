<script context="module">
	const privateRegex = /:private/;

	function getTags(todo, text) {
		const tags = {};

		if (todo && todo.due_date) {
			tags['due_date'] = { value: new Date(todo.due_date) };
		}

		if (todo && todo.private) {
			tags['private'] = { value: 'private' };
		}

		const dates = chrono.parse(text, new Date(), { forwardDate: true });
		if (dates.length !== 0) {
			const date = dates[0];
			tags['due_date'] = {
				value: date.date(),
				start: date.index,
				end: date.index + date.text.length
			};
		}

		const privateMatches = text.match(privateRegex);
		if (privateMatches) {
			tags['private'] = {
				value: 'private',
				start: privateMatches.index,
				end: privateMatches.index + privateMatches[0].length
			};
		}

		return tags;
	}
</script>

<script>
	import { api } from '$lib/api';
	import { session } from '$app/stores';
	import * as chrono from 'chrono-node';
	import Edit, { getStripped } from './Edit.svelte';

	export let edit = undefined;
	export let todos = undefined;
	export let todo = undefined;

	let text = todo !== undefined ? todo.description : '';

	const fg = {
		due_date: 'white',
		private: 'white'
	};

	const bg = {
		due_date: 'var(--blue-dark)',
		private: 'var(--turquoise-dark)'
	};

	async function todoCreate(e) {
		e.preventDefault();
		if (text.length === 0) return;

		const tags = getTags(todo, text);

		const todoDiff = {
			description: getStripped(text, tags),
			due_date: tags['due_date'] ? tags['due_date'].value : undefined,
			private: tags['private'] !== undefined
		};

		if (todo !== undefined) {
			Object.assign(todo, todoDiff);
			todo = todo;
			await api('PUT', fetch, $session, `todo/${todo.id}`, todo);
			edit = undefined;
		} else {
			const newTodo = await api('POST', fetch, $session, `todo`, todoDiff);
			todos = [...todos, newTodo];
			text = '';
		}
	}

	$: tags = getTags(todo, text);
</script>

<Edit
	on:submit={todoCreate}
	bind:text
	placeholder="Todo description"
	cancel={todo ? () => (edit = undefined) : undefined}
	{tags}
	{fg}
	{bg}
/>

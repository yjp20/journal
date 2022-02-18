<script context="module">
	const blockedRegex = /[!]?blocked/
	const privateRegex = /[!]?private/
	const recurRegex = /[!]?recur(?::(\d+))?/

	function getTags(todo, text) {
		const tags = {}

		if (todo && todo.due_date) {
			const date = new Date(todo.due_date)
			tags['due_date'] = { display: date, value: date }
		}

		if (todo && todo.private) {
			tags['private'] = { value: true }
		}

		if (todo && todo.blocked) {
			tags['blocked'] = { value: true }
		}

		const parsedDates = chrono.parse(text, new Date(), { forwardDate: true })
		if (parsedDates.length !== 0) {
			const parsedDate = parsedDates[0]
			const component = parsedDate.start
			component.impliedValues['hour'] = 23
			component.impliedValues['minute'] = 59
			component.impliedValues['second'] = 0
			component.impliedValues['milisecond'] = 0

			tags['due_date'] = {
				value: component.date(),
				display: component.date(),
				start: parsedDate.index,
				end: parsedDate.index + parsedDate.text.length
			}
		}

		const blockedMatches = text.match(blockedRegex)
		if (blockedMatches) {
			tags['blocked'] = {
				value: blockedMatches[0][0] != '!',
				start: blockedMatches.index,
				end: blockedMatches.index + blockedMatches[0].length
			}
		}

		const privateMatches = text.match(privateRegex)
		if (privateMatches) {
			tags['private'] = {
				value: privateMatches[0][0] != '!',
				start: privateMatches.index,
				end: privateMatches.index + privateMatches[0].length
			}
		}

		const recurMatches = text.match(recurRegex)
		if (recurMatches) {
			tags['recur'] = {
				value: parseInt(recurMatches[1]),
				display: `repeat after ${recurMatches[1]} days`,
				start: recurMatches.index,
				end: recurMatches.index + recurMatches[0].length
			}
		}

		return tags
	}
</script>

<script>
	import { createEventDispatcher } from 'svelte'
	import * as chrono from 'chrono-node'

	import { api } from '$lib/api'
	import { session } from '$app/stores'
	import Edit, { getStripped } from './Edit.svelte'

	export let todos
	export let todo = undefined

	let text = todo !== undefined ? todo.description : ''

	const fg = {
		due_date: 'white',
		private: 'white'
	}

	const bg = {
		due_date: 'var(--blue-dark)',
		private: 'var(--turquoise-dark)'
	}

	const dispatch = createEventDispatcher()

	async function todoCreate(e) {
		e.preventDefault()
		if (text.length === 0) return

		const tags = getTags(todo, text)

		const todoDiff = {
			description: getStripped(text, tags),
			due_date: tags['due_date'] ? tags['due_date'].value : undefined,
			private: !!tags['private']?.value,
			blocked: !!tags['blocked']?.value,
			recur: tags['recur']?.value == 0 ? undefined : tags['recur']?.value
		}

		if (todo !== undefined) {
			Object.assign(todo, todoDiff)
			todo = todo
			if (todo.recur && !todo.due_date) {
				todo.recur = undefined
			}
			await api('PUT', fetch, $session, `todo/${todo.id}`, todo)
			dispatch('finish')
		} else {
			if (todoDiff.recur && !todoDiff.due_date) {
				alert('Cannot recur without a due date')
				return
			}
			const newTodo = await api('POST', fetch, $session, `todo`, todoDiff)
			todos = [...todos, newTodo]
			text = ''
		}
	}

	$: tagsPromise = getTags(todo, text)
</script>

<Edit
	on:submit={todoCreate}
	bind:text
	placeholder="Todo description"
	on:cancel={() => dispatch('finish')}
	cancel={!!todo}
	{tagsPromise}
	{fg}
	{bg} />

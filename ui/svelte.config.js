import preprocess from 'svelte-preprocess';
import node from '@sveltejs/adapter-node';

const config = {
	preprocess: preprocess(),
	kit: {
		adapter: node(),
		target: '#svelte'
	}
};

export default config;

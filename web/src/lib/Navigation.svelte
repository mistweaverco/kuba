<script lang="ts">
	import { page } from '$app/stores';

	$: currentPath = $page.url.pathname;

	const navItems = [
		{ href: '/', label: 'Home' },
		{ href: '/installation', label: 'Installation' },
		{ href: '/usage', label: 'Usage' },
		{ href: '/configuration', label: 'Configuration' },
		{ href: '/providers', label: 'Providers' },
		{ href: '/examples', label: 'Examples' }
	];
</script>

<nav class="navbar bg-base-100 shadow-lg">
	<div class="navbar-start">
		<div class="dropdown">
			<div tabindex="0" role="button" class="btn btn-ghost lg:hidden">
				<i class="fa-solid fa-bars"></i>
			</div>
			<ul class="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 rounded-box w-52">
				{#each navItems as item}
					<li>
						<a href={item.href} class={currentPath === item.href ? 'active' : ''}>
							{item.label}
						</a>
					</li>
				{/each}
			</ul>
		</div>
		<a href="/" class="btn btn-ghost text-xl">
			<img src="/logo.svg" alt="Kuba" class="w-8 h-8 mr-2" />
			Kuba
		</a>
	</div>
	<div class="navbar-center hidden lg:flex">
		<ul class="menu menu-horizontal px-1">
			{#each navItems as item}
				<li>
					<a href={item.href} class={currentPath === item.href ? 'active' : ''}>
						{item.label}
					</a>
				</li>
			{/each}
		</ul>
	</div>
	<div class="navbar-end">
		<a
			href="https://github.com/mistweaverco/kuba"
			class="btn btn-ghost btn-circle"
			aria-label="View Kuba on GitHub"
		>
			<i class="fa-brands fa-github fa-lg"></i>
		</a>
	</div>
</nav>

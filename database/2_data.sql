INSERT INTO users (id, password, name, email, role, balance)
VALUES (
		'460aec7e-7d58-42fd-93b8-bca05a77bbf5',
		'$2a$10$L22SyyVSBP3WZAm7K/bfEOPfMCZsQhQEPs9vwRfbSAd.MIcoCp6rC',
		'John',
		'john@example.com',
		1,
		10
	),
	(
		'8c3524e5-a297-42aa-85d3-faca261cbfb8',
		'$2a$10$L22SyyVSBP3WZAm7K/bfEOPfMCZsQhQEPs9vwRfbSAd.MIcoCp6rC',
		'Marc',
		'marc@example.com',
		0,
		10
	),
	(
		'3b4fef91-2523-46ab-b06d-17e3e2d4b209',
		'$2a$10$L22SyyVSBP3WZAm7K/bfEOPfMCZsQhQEPs9vwRfbSAd.MIcoCp6rC',
		'Jason',
		'jason@example.com',
		0,
		10
	),
	(
		'80ddee0a-b1cc-4c03-8a78-b994486850e7',
		'$2a$10$L22SyyVSBP3WZAm7K/bfEOPfMCZsQhQEPs9vwRfbSAd.MIcoCp6rC',
		'Martin',
		'martin@example.com',
		0,
		10
	);
INSERT INTO promotions (
		id,
		title,
		description,
		amount,
		is_active,
		type
	)
VALUES (
		'abfe0ac0-1c51-4ca5-974d-13dcb1c76d43',
		'Welcome bonus',
		'Bonus that is earned by register',
		20,
		true,
		'welcome_bonus'
	),
	(
		'ad0b6132-0f06-49d2-bdfd-4d767ac5b0fc',
		'Cool bonus',
		'Bonus that is earned by being cool',
		20,
		true,
		'regular'
	);
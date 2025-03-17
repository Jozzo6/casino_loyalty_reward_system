INSERT INTO promotions (
	id,
	title,
	description,
	amount,
	is_active,
	type
) VALUES (
	gen_random_uuid(),
	'Welcome bonus',
	'Bonus that is earned by register',
	20,
	true,
	'welcome_bonus'
)
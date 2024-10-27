Select name to name relations:

```sql
SELECT
	parent.name,
	child.name
FROM
	article AS parent
	JOIN relation ON parent.id = from_id
	INNER JOIN article AS child ON child.id = to_id
WHERE
	parent.status = 'COMPLETED'
    AND child.status = 'COMPLETED'
```

TODO: backlinks

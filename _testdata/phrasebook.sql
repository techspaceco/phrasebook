-- wu
-- tang

-- woot
-- export: TestExportYesql, {"foo":"bar","baz":1}
select *
from test_export_yesql
where id = @id;
-- end

-- test
-- baz
-- export: TestExportPhrasebook
select
  *
from
  test_export_phrasebook
where
  id = @id;
-- end

-- export TestExportCGO
select * from test_export_cgo where "id" = @id;
-- end

--export TestExportNoWhitespaceCGO
select * from test_export_no_whitesapce_cgo where id = @id;

-- export TestExportJSON {"foo":1}
select * from test_export_json where id = @id;

-- leading
-- comment
-- block
--
-- name: TestLeadingCommentBlock
select * from test_leading_comment_block where id = @id;

-- leading
-- compact
-- comment
-- block
-- name: TestLeadingCompactCommentBlock
select * from test_leading_compact_comment_block where id = @id;

--no
--leading
--whitespace
--comment
--block
-- name: TestNoLeadingWhitespaceCommentBlock
select * from test_no_leading_whitespace_comment_block where id = @id;
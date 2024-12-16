CREATE OR REPLACE FUNCTION search_products(
    in_branch_id VARCHAR(15),
    in_keywords TEXT
)
RETURNS TABLE (
    product_id VARCHAR(19),
    product_name VARCHAR(100),
    sales_price BIGINT,
    alternate_price BIGINT,
    stock BIGINT,
    unit_id VARCHAR(15),
    unit_name VARCHAR(100)
)
LANGUAGE plpgsql
AS $$
DECLARE
    keyword_array TEXT[];
    keyword_count INT;
BEGIN
    -- Pisahkan keywords menjadi array
    keyword_array := string_to_array(in_keywords, ' ');
    keyword_count := array_length(keyword_array, 1); -- Hitung jumlah keyword

    -- Jalankan query dengan validasi semua keyword menggunakan trigram similarity
    RETURN QUERY
    SELECT 
        p.id AS product_id,
        p.name AS product_name,
        p.sales_price AS sales_price,
        p.alternate_price AS alternate_price,
        p.stock,
        p.unit_id AS unit_id,
        u.name AS unit_name
    FROM products p
    JOIN units u ON p.unit_id = u.id
    WHERE p.branch_id = in_branch_id
      AND p.stock > 0
      AND (
          SELECT COUNT(*)
          FROM unnest(keyword_array) AS input_word
          WHERE (
              -- Cek kemiripan nama produk
              similarity(input_word, p.name) > 0.3 OR
              -- Cek kemiripan deskripsi produk
              similarity(input_word, p.description) > 0.3
          )
      ) = keyword_count
    ORDER BY p.name ASC;
END;
$$;
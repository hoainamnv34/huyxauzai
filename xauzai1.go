def insert_content(file_a, file_b, position='end'):
    """
    Chèn nội dung của file_a vào file_b.

    :param file_a: Đường dẫn file nguồn (file a).
    :param file_b: Đường dẫn file đích (file b).
    :param position: Vị trí chèn ('start' hoặc 'end'). Mặc định là 'end'.
    """
    try:
        # Đọc nội dung file a
        with open(file_a, 'r') as fa:
            content_a = fa.read()

        # Đọc nội dung file b
        with open(file_b, 'r') as fb:
            content_b = fb.read()

        # Chèn nội dung
        if position == 'start':
            new_content = content_a + '\n' + content_b
        elif position == 'end':
            new_content = content_b + '\n' + content_a
        else:
            raise ValueError("Invalid position. Use 'start' or 'end'.")

        # Ghi nội dung mới vào file b
        with open(file_b, 'w') as fb:
            fb.write(new_content)

        print(f"Successfully inserted content from {file_a} into {file_b} at {position}.")
    except Exception as e:
        print(f"Failed to insert content from {file_a} into {file_b}. Error: {e}")

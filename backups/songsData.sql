PGDMP                      }            songs_db    16.0    16.0 	    �           0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
                      false            �           0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                      false            �           0    0 
   SEARCHPATH 
   SEARCHPATH     8   SELECT pg_catalog.set_config('search_path', '', false);
                      false            �           1262    50342    songs_db    DATABASE     �   CREATE DATABASE songs_db WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'English_United States.1252';
    DROP DATABASE songs_db;
                postgres    false            �          0    50351    groups 
   TABLE DATA           *   COPY public.groups (id, name) FROM stdin;
    public          postgres    false    217   2       �          0    50343    schema_migrations 
   TABLE DATA           ;   COPY public.schema_migrations (version, dirty) FROM stdin;
    public          postgres    false    215   �       �          0    50358    songs 
   TABLE DATA           O   COPY public.songs (id, "groupId", name, "releaseDate", text, link) FROM stdin;
    public          postgres    false    219          �           0    0    groups_id_seq    SEQUENCE SET     <   SELECT pg_catalog.setval('public.groups_id_seq', 17, true);
          public          postgres    false    216            �           0    0    songs_id_seq    SEQUENCE SET     :   SELECT pg_catalog.setval('public.songs_id_seq', 9, true);
          public          postgres    false    218            �   �   x�%���@�����O �(�J��� ���&����̴���;9����z���FZV�5c@�N�aH���aD�p��PЎ�M��G��oLh-��"!�Z�)�r�>mm|�cC�v&��HȪ.�Uʵ�<̐��d�.~{���v�o�������-տ�� | @	5�      �      x�3�L����� �T      �   �  x��Xko�6���
~�8m�5i:`(�>V�ڡ-�Z�-�������Rvm#]�k:`@�����{xo�zǽ����R����\�������������W��Vx�sq��A4�,D(��S;��J䲱�����"L앆�F#��5B��D떵��LM���ͭ������*_�U��r:SbV�R�JU��u�q8Lҙ�z��z���S��t|�Z兄T]����}sr��[��ӹ�K�c!K��<xS�g�J��A��|QdlS4h�3R��j;>Ț�n/L�W&ƩS�������rb�*�����1��f�c^~���jd����ᐵ���	A�V2+bJf�N��1��!�9r�� WP�Gӛ��V
0���2�w)�w�����x�^B���i�{��Su/s��F��x��������I}�<�=�BJG��s��W�-x�x����y��0�k��[z �`Z<��� Ⱥ�c�q�(��-�C���a�/���J��MV���!� 4cP�lu,z�/݂x$J�њV����^ ���V��r�I�eU�a%�i���E�P��a�%���λb��g�z�����)]��L�RZ�`:�Pr��F/p�¹(-��!>���4��v�Tp	���abʧrڊ��;� ځ0ʣk�a��E�F�t(�#1�]c��U�+��49"�C�ʂ^�6z-fJ���	o3E��*ӜLr�J����2h?�*
]1SU�X(��`�Ĝ�Ӎ�/d�����Ɉ��  Ԭ����L��G�kE�s1�LD���G&�n�e���V!��g����X��#�k�X�K���&R1n��Jt�"9��e!��Q�PXl�O��S	�_B����[���0�d��;J�w7�m�ʫ'^=�8}x�:��{'��@	�8��&���������;�a��dFψ
,	�X��!����Z�F��(p���$�?@k��7m�#��'�1h^B����nɒ���g[
s�|9����ci�>��!}�c�~Ƈ�-�Dg��|"]d!��t���������#v&A��Qb�����6�㊟q%6�l����}G��YD�3�+�_��,�3�ސO�P&��6r[Y�S#KV.Q,���\�(�8vb��s����H��;^�{�s@�z��&S(������I�#;��,�����`����r��xg9�g��2(�~o��[�����F�9��"
�,~w�G�����u�e^����5�9�������E�!��5�Ę����v��w�^��S�}IC�P�*��O�m쫷��������cv��Nq��
���)o�gh�-��[��.���1`��,*m��]ciK���6Ut,�Ĥ��1�1y���w#Eb{�8$S 
uqLś4��ы�i`�"[L
(�9�ҿ�eW<�Sλ���IO�Q�������sj~��U�U��p�������������It�&���`QY�M�4K*�SGq�a���u�U�m���.��	�w��b;B�
}�G����]��<{�E��/[�6W8I~4���JY����J�t%���@����}4&)ilc̺����lF X�U�d= ��At�c"`k� Q�cq��:��8���Gz7�s�j`S���
�ńoxp��zW�ǟ?]��&����������V~�#     
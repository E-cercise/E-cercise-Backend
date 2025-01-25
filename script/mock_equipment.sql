INSERT INTO equipment (id, name, price, brand, model, color, material, weight, remaining_products, special_feature)
VALUES
    (uuid_generate_v4(), 'Back Muscle Trainer', 199.99, 'Brand A', 'Model X1', 'Black', 'Steel', 20.5, 10, 'Targets back muscles'),
    (uuid_generate_v4(), 'Thigh Strengthener', 149.99, 'Brand B', 'Model Y2', 'Gray', 'Aluminum', 15.0, 5, 'Focuses on thigh and hamstring muscles'),
    (uuid_generate_v4(), 'Calf Stretcher', 99.99, 'Brand C', 'Model Z3', 'Blue', 'Plastic', 5.0, 20, 'Helps improve calf flexibility');


INSERT INTO equipment_muscle_groups (equipment_id, muscle_group_id)
VALUES
-- Associate 'Back Muscle Trainer' with back muscles
((SELECT id FROM equipment WHERE name = 'Back Muscle Trainer'), 'bk_1'),
((SELECT id FROM equipment WHERE name = 'Back Muscle Trainer'), 'bk_2'),
((SELECT id FROM equipment WHERE name = 'Back Muscle Trainer'), 'bk_3'),
-- Associate 'Thigh Strengthener' with thigh muscles
((SELECT id FROM equipment WHERE name = 'Thigh Strengthener'), 'bk_20'),
((SELECT id FROM equipment WHERE name = 'Thigh Strengthener'), 'ft_24'),
-- Associate 'Calf Stretcher' with calf muscles
((SELECT id FROM equipment WHERE name = 'Calf Stretcher'), 'bk_30'),
((SELECT id FROM equipment WHERE name = 'Calf Stretcher'), 'ft_38');
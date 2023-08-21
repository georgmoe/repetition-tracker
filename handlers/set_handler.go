package handlers

// // receive a set and append to existing exercise
// func PostSet(c *fiber.Ctx) error {
// 	// get primitive user id
// 	userIdFromLocals := c.Locals(USER_ID)
// 	userId, err := GetPrimitiveObjectIDFromInterface(userIdFromLocals)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
// 	}

// 	// get workout id from path
// 	workoutIdStr := c.Params("workoutId")
// 	workoutId, err := primitive.ObjectIDFromHex(workoutIdStr)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
// 	}
// }

// func Put(c *fiber.Ctx) error {

// }

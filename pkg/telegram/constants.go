package telegram

const (
	startCmd             = "/start"
	helpCmd              = "/help"
	addCategoryCmd       = "/add_category"
	getCategoriesListCmd = "/my_categories"
	addExpensesCmd       = "/add_expenses"
	backCmd              = "/back"
	deleteCategoryCmd    = "/delete_category"
	deleteAllCmd         = "/delete_all"

	msgYes               = "ДА"
	msgNo                = "НЕТ"
	msgConfirmation      = "Вы уверены, что хотите удалить все категории?\nВведите «ДА/НЕТ»:"
	msgEnterCatrgoryData = "Введите название категории и лимит трат. Например: «Название категории 1000».\nЧтобы отменить действие, используйте команду /back\n"
	msgEnterExpensesData = "Введите название категории и сумму трат, которую хотите добавить. Например: «Название категории 1000».\nЧтобы отменить действие, используйте команду /back\n"
	msgCategoriesList    = "Список Ваших категорий\n\n"
	msgStart             = "Привет, я Кэш Кипер!👾\nЯ помогу тебе с финансами. Чтобы узнать о работе бота, отправьте команду /help \n"
	msgHelp              = "Список команд:\n/add_category - создать категорию\n/my_categories - получить список ваших категорий\n/delete_category - удалить категорию\n/add_expenses - добавить траты\n"
	msgSuccess           = "Успешно! "
	msgCategoryDelete    = "Введите название категории для удаления. Чтобы отменить действие, используйте команду /back\n"
	msgEnterCategoryName = "Введите название категории, в которую хотите добавить траты. Чтобы отменить действие, используйте команду /back\n"

	errMsgTimout             = "Превышено время ожидания"
	errMsgUnknownCmd         = "Неизвестная команда. "
	errMsgDefault            = "Ошибка. "
	errMsgWrongTitle         = "Название не может быть командой! "
	errMsgUnknownTitle       = "Категория не найдена. "
	errMsgDeny               = "Действие отменено. "
	errMsgWrongMassageFormat = "Неверный формат создания категории. "
	errMsgWrongTotalValue    = "Сумма трат не может быть 0 или меньше. "
	errMsgAlreadyUsed        = "Такое название категории уже существует. "

	minTotalValue     = 1
	minUserMessageLen = 2
)

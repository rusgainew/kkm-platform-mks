import UserCreateForm from '@/features/users/components/UserCreateForm';

export const metadata = {
  title: 'Создание пользователя | KKM',
  description: 'Форма создания нового пользователя',
};

export default function CreateUserPage() {
  return (
    <div className="min-h-screen bg-gray-950 p-4 md:p-8">
      <UserCreateForm />
    </div>
  );
}

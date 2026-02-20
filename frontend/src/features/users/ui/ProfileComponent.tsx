import { CopyOutlined, EyeInvisibleOutlined, EyeTwoTone, LoadingOutlined } from '@ant-design/icons';
import { App, Button, Input, Modal, Spin, Typography } from 'antd';
import { useEffect, useState } from 'react';

import { apiKeyApi, userApi } from '../../../entity/users';
import type { ApiKeyInfo } from '../../../entity/users/model/ApiKeyInfo';
import type { ChangePasswordRequest } from '../../../entity/users/model/ChangePasswordRequest';
import type { CreateApiKeyResponse } from '../../../entity/users/model/CreateApiKeyResponse';
import type { SignInRequest } from '../../../entity/users/model/SignInRequest';
import type { UpdateUserInfoRequest } from '../../../entity/users/model/UpdateUserInfoRequest';
import type { UserProfile } from '../../../entity/users/model/UserProfile';
import { UserRole } from '../../../entity/users/model/UserRole';

const { Text } = Typography;

interface Props {
  contentHeight: number;
}

const getRoleDisplayText = (role: UserRole): string => {
  switch (role) {
    case UserRole.ADMIN:
      return 'Admin';
    case UserRole.MEMBER:
      return 'Member';
    default:
      return role;
  }
};

const formatDate = (dateString: string | null): string => {
  if (!dateString) return 'Never';
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
};

export function ProfileComponent({ contentHeight }: Props) {
  const { message } = App.useApp();
  const [user, setUser] = useState<UserProfile | undefined>(undefined);
  const [isChangingPassword, setIsChangingPassword] = useState(false);

  const [editName, setEditName] = useState('');
  const [editEmail, setEditEmail] = useState('');
  const [isUpdatingProfile, setIsUpdatingProfile] = useState(false);
  const [editNameError, setEditNameError] = useState(false);
  const [editEmailError, setEditEmailError] = useState(false);

  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [newPasswordVisible, setNewPasswordVisible] = useState(false);
  const [confirmPasswordVisible, setConfirmPasswordVisible] = useState(false);

  const [newPasswordError, setNewPasswordError] = useState(false);
  const [confirmPasswordError, setConfirmPasswordError] = useState(false);

  const [apiKeyInfo, setApiKeyInfo] = useState<ApiKeyInfo | null>(null);
  const [isLoadingApiKey, setIsLoadingApiKey] = useState(false);
  const [showApiKeyModal, setShowApiKeyModal] = useState(false);
  const [generatedApiKey, setGeneratedApiKey] = useState<string | null>(null);
  const [isGeneratingKey, setIsGeneratingKey] = useState(false);

  useEffect(() => {
    loadUserProfile();
    loadApiKeyInfo();
  }, []);

  const loadUserProfile = () => {
    userApi
      .getCurrentUser()
      .then((user) => {
        setUser(user);
        setEditName(user.name);
        setEditEmail(user.email);
      })
      .catch((error) => {
        message.error(error.message);
      });
  };

  const loadApiKeyInfo = () => {
    setIsLoadingApiKey(true);
    apiKeyApi
      .getApiKey()
      .then((info) => {
        setApiKeyInfo(info);
      })
      .catch(() => {
        setApiKeyInfo(null);
      })
      .finally(() => {
        setIsLoadingApiKey(false);
      });
  };

  const validatePasswordFields = (): boolean => {
    let isValid = true;

    if (!newPassword) {
      setNewPasswordError(true);
      isValid = false;
    } else if (newPassword.length < 6) {
      setNewPasswordError(true);
      message.error('Password must be at least 6 characters long');
      isValid = false;
    } else {
      setNewPasswordError(false);
    }

    if (!confirmPassword) {
      setConfirmPasswordError(true);
      isValid = false;
    } else if (newPassword !== confirmPassword) {
      setConfirmPasswordError(true);
      message.error('New passwords do not match');
      isValid = false;
    } else {
      setConfirmPasswordError(false);
    }

    return isValid;
  };

  const handlePasswordChange = async () => {
    if (!validatePasswordFields()) {
      return;
    }

    setIsChangingPassword(true);

    try {
      const request: ChangePasswordRequest = {
        newPassword,
      };

      await userApi.changePassword(request);

      setNewPassword('');
      setConfirmPassword('');

      if (user?.email) {
        try {
          const signInRequest: SignInRequest = {
            email: user.email,
            password: newPassword,
          };
          await userApi.signIn(signInRequest);
          message.success('Successfully signed in with new password');
        } catch (signInError: unknown) {
          const errorMessage =
            signInError instanceof Error
              ? signInError.message
              : 'Failed to sign in with new password';
          message.error(errorMessage);
          userApi.logout();
          userApi.notifyAuthListeners();
          window.location.reload();
        }
      }
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to change password';
      message.error(errorMessage);
    } finally {
      setIsChangingPassword(false);
    }
  };

  const handleProfileUpdate = async () => {
    if (!editName || editName.trim() === '') {
      setEditNameError(true);
      message.error('Name is required');
      return;
    }
    setEditNameError(false);

    if (user?.email !== 'admin') {
      if (!editEmail || editEmail.trim() === '') {
        setEditEmailError(true);
        message.error('Email is required');
        return;
      }
      setEditEmailError(false);
    }

    setIsUpdatingProfile(true);

    try {
      const request: UpdateUserInfoRequest = {};

      if (editName !== user?.name) {
        request.name = editName;
      }
      if (user?.email !== 'admin' && editEmail !== user?.email) {
        request.email = editEmail;
      }

      if (Object.keys(request).length === 0) {
        message.info('No changes to save');
        setIsUpdatingProfile(false);
        return;
      }

      await userApi.updateUserInfo(request);
      message.success('Profile updated successfully');

      loadUserProfile();
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to update profile';
      message.error(errorMessage);
    } finally {
      setIsUpdatingProfile(false);
    }
  };

  const handleLogout = () => {
    userApi.logout();
    window.location.reload();
  };

  const handleGenerateApiKey = async () => {
    setIsGeneratingKey(true);

    try {
      let response: CreateApiKeyResponse;

      if (apiKeyInfo) {
        response = await apiKeyApi.regenerateApiKey();
      } else {
        response = await apiKeyApi.createApiKey();
      }

      setGeneratedApiKey(response.apiKey);
      setShowApiKeyModal(true);
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to generate API key';
      message.error(errorMessage);
    } finally {
      setIsGeneratingKey(false);
    }
  };

  const handleCopyApiKey = () => {
    if (generatedApiKey) {
      navigator.clipboard.writeText(generatedApiKey);
      message.success('API key copied to clipboard');
    }
  };

  const handleModalClose = () => {
    setShowApiKeyModal(false);
    setGeneratedApiKey(null);
    loadApiKeyInfo();
  };

  return (
    <div className="flex grow">
      <div className="w-full">
        <div
          className="grow overflow-y-auto rounded bg-white p-5 shadow dark:bg-gray-800"
          style={{ height: contentHeight }}
        >
          <h1 className="text-2xl font-bold dark:text-white">Profile</h1>

          <div className="mt-5">
            {user ? (
              <>
                <div className="mb-6">
                  <h3 className="mb-4 text-lg font-semibold dark:text-white">
                    Profile Information
                  </h3>
                  <div className="max-w-md">
                    <div className="text-xs font-semibold dark:text-gray-200">User ID</div>
                    <div className="mb-4 text-sm text-gray-600 dark:text-gray-400">{user.id}</div>

                    <div className="mb-1 text-xs font-semibold dark:text-gray-200">Name</div>
                    <Input
                      value={editName}
                      onChange={(e) => {
                        setEditNameError(false);
                        setEditName(e.currentTarget.value);
                      }}
                      status={editNameError ? 'error' : undefined}
                      placeholder="Enter your name"
                      className="mb-4"
                    />

                    <div className="mt-2 mb-1 text-xs font-semibold dark:text-gray-200">Email</div>
                    <Input
                      value={editEmail}
                      onChange={(e) => {
                        setEditEmailError(false);
                        setEditEmail(e.currentTarget.value.trim().toLowerCase());
                      }}
                      status={editEmailError ? 'error' : undefined}
                      placeholder="Enter your email"
                      type="email"
                      className="mb-4"
                      disabled={user.email === 'admin'}
                    />
                    {user.email === 'admin' && (
                      <div className="mb-4 text-xs text-gray-500 dark:text-gray-400">
                        Admin email cannot be changed
                      </div>
                    )}

                    <div className="mt-2 mb-1 text-xs font-semibold dark:text-gray-200">Role</div>
                    <div className="mb-4">
                      <span className="inline-flex items-center rounded-full bg-blue-100 px-2.5 py-0.5 text-xs font-medium text-blue-800 dark:bg-blue-900 dark:text-blue-200">
                        {getRoleDisplayText(user.role)}
                      </span>
                    </div>

                    {(editName !== user.name || editEmail !== user.email) && (
                      <Button
                        type="primary"
                        onClick={handleProfileUpdate}
                        loading={isUpdatingProfile}
                        disabled={isUpdatingProfile}
                        className="border-blue-600 bg-blue-600 hover:border-blue-700 hover:bg-blue-700"
                      >
                        Save changes
                      </Button>
                    )}
                  </div>
                </div>

                <div className="mb-8">
                  <Button type="primary" ghost onClick={handleLogout} danger>
                    Logout
                  </Button>
                </div>

                <div className="max-w-xs">
                  <h3 className="mb-4 text-lg font-semibold dark:text-white">Change Password</h3>

                  <div className="max-w-sm">
                    <div className="my-1 text-xs font-semibold dark:text-gray-200">
                      New Password
                    </div>
                    <Input.Password
                      placeholder="Enter new password"
                      value={newPassword}
                      onChange={(e) => {
                        setNewPasswordError(false);
                        setNewPassword(e.currentTarget.value);
                      }}
                      status={newPasswordError ? 'error' : undefined}
                      iconRender={(visible) =>
                        visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />
                      }
                      visibilityToggle={{
                        visible: newPasswordVisible,
                        onVisibleChange: setNewPasswordVisible,
                      }}
                      autoComplete="new-password"
                    />

                    <div className="mt-2 mb-1 text-xs font-semibold dark:text-gray-200">
                      Confirm New Password
                    </div>
                    <Input.Password
                      placeholder="Confirm new password"
                      value={confirmPassword}
                      onChange={(e) => {
                        setConfirmPasswordError(false);
                        setConfirmPassword(e.currentTarget.value);
                      }}
                      status={confirmPasswordError ? 'error' : undefined}
                      iconRender={(visible) =>
                        visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />
                      }
                      autoComplete="new-password"
                      visibilityToggle={{
                        visible: confirmPasswordVisible,
                        onVisibleChange: setConfirmPasswordVisible,
                      }}
                    />

                    <div className="mt-3" />

                    {(newPassword || confirmPassword) && (
                      <Button
                        type="primary"
                        onClick={handlePasswordChange}
                        loading={isChangingPassword}
                        disabled={isChangingPassword}
                        className="border-blue-600 bg-blue-600 hover:border-blue-700 hover:bg-blue-700"
                      >
                        {isChangingPassword ? 'Changing password...' : 'Change password'}
                      </Button>
                    )}
                  </div>
                </div>

                <div className="mt-8 max-w-md">
                  <h3 className="mb-4 text-lg font-semibold dark:text-white">API Key</h3>

                  {isLoadingApiKey ? (
                    <Spin indicator={<LoadingOutlined spin />} />
                  ) : apiKeyInfo ? (
                    <div>
                      <div className="mb-1 text-xs font-semibold dark:text-gray-200">Key</div>
                      <div className="mb-4 rounded bg-gray-100 p-2 font-mono text-sm dark:bg-gray-700 dark:text-gray-200">
                        {apiKeyInfo.keyPrefix}...
                      </div>

                      <div className="mb-1 text-xs font-semibold dark:text-gray-200">Created</div>
                      <div className="mb-4 text-sm text-gray-600 dark:text-gray-400">
                        {formatDate(apiKeyInfo.createdAt)}
                      </div>

                      <div className="mb-1 text-xs font-semibold dark:text-gray-200">Last Used</div>
                      <div className="mb-4 text-sm text-gray-600 dark:text-gray-400">
                        {formatDate(apiKeyInfo.lastUsedAt)}
                      </div>

                      <Button
                        type="primary"
                        onClick={handleGenerateApiKey}
                        loading={isGeneratingKey}
                        disabled={isGeneratingKey}
                        danger
                      >
                        Regenerate API Key
                      </Button>
                    </div>
                  ) : (
                    <div>
                      <p className="mb-4 text-sm text-gray-600 dark:text-gray-400">
                        No API key configured. Generate one to access the API.
                      </p>
                      <Button
                        type="primary"
                        onClick={handleGenerateApiKey}
                        loading={isGeneratingKey}
                        disabled={isGeneratingKey}
                        className="border-blue-600 bg-blue-600 hover:border-blue-700 hover:bg-blue-700"
                      >
                        Generate API Key
                      </Button>
                    </div>
                  )}
                </div>
              </>
            ) : (
              <div>
                <Spin indicator={<LoadingOutlined spin />} />
              </div>
            )}
          </div>
        </div>
      </div>

      <Modal
        title="API Key Generated"
        open={showApiKeyModal}
        onOk={handleModalClose}
        onCancel={handleModalClose}
        footer={[
          <Button key="close" type="primary" onClick={handleModalClose}>
            Close
          </Button>,
        ]}
      >
        <div className="mb-4">
          <Text type="warning" className="mb-2 block font-semibold">
            Make sure to copy your API key now. You won't be able to see it again!
          </Text>
        </div>

        <div className="mb-4 rounded bg-gray-100 p-3 dark:bg-gray-700">
          <Text copyable={{ text: generatedApiKey || '', onCopy: handleCopyApiKey }}>
            <code className="text-sm">{generatedApiKey}</code>
          </Text>
        </div>

        <Button icon={<CopyOutlined />} onClick={handleCopyApiKey} className="w-full">
          Copy to Clipboard
        </Button>
      </Modal>
    </div>
  );
}

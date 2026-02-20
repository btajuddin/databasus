import { CopyOutlined, LoadingOutlined } from '@ant-design/icons';
import { App, Button, Modal, Spin, Typography } from 'antd';
import { useEffect, useState } from 'react';

import { apiKeyApi } from '../../../entity/users';
import type { ApiKeyInfo } from '../../../entity/users/model/ApiKeyInfo';

const { Text } = Typography;

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

export function ApiKeyComponent() {
  const { message } = App.useApp();

  const [apiKeyInfo, setApiKeyInfo] = useState<ApiKeyInfo | null>(null);
  const [isLoadingApiKey, setIsLoadingApiKey] = useState(false);
  const [showApiKeyModal, setShowApiKeyModal] = useState(false);
  const [generatedApiKey, setGeneratedApiKey] = useState<string | null>(null);
  const [isGeneratingKey, setIsGeneratingKey] = useState(false);

  useEffect(() => {
    loadApiKeyInfo();
  }, []);

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

  const handleGenerateApiKey = async () => {
    setIsGeneratingKey(true);

    try {
      const response = await apiKeyApi.upsertApiKey();
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
    <div className="max-w-md">
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

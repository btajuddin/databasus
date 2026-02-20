import { getApplicationServer } from '../../../constants';
import RequestOptions from '../../../shared/api/RequestOptions';
import { apiHelper } from '../../../shared/api/apiHelper';
import type { ApiKeyInfo } from '../model/ApiKeyInfo';
import type { CreateApiKeyResponse } from '../model/CreateApiKeyResponse';

export const apiKeyApi = {
  async getApiKey(): Promise<ApiKeyInfo> {
    const requestOptions = new RequestOptions();
    return apiHelper.fetchGetJson(
      `${getApplicationServer()}/api/v1/users/me/api-key`,
      requestOptions,
    );
  },

  async createApiKey(): Promise<CreateApiKeyResponse> {
    const requestOptions = new RequestOptions();
    return apiHelper.fetchPostJson(
      `${getApplicationServer()}/api/v1/users/me/api-key`,
      requestOptions,
    );
  },

  async regenerateApiKey(): Promise<CreateApiKeyResponse> {
    const requestOptions = new RequestOptions();
    return apiHelper.fetchPostJson(
      `${getApplicationServer()}/api/v1/users/me/api-key/regenerate`,
      requestOptions,
    );
  },
};

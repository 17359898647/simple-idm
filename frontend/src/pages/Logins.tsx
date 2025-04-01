import { useNavigate } from '@solidjs/router';
import type { Component } from 'solid-js';
import { createSignal, For, onMount, Show } from 'solid-js';

import type { Login } from '../api/login';
import { loginApi } from '../api/login';

const Logins: Component = () => {
  const navigate = useNavigate();
  const [logins, setLogins] = createSignal<Login[]>([]);
  const [error, setError] = createSignal<string | null>(null);
  const [loading, setLoading] = createSignal(true);

  const fetchLogins = async () => {
    try {
      const data = await loginApi.listLogins();
      setLogins(data);
    }
    catch (err) {
      setError(err instanceof Error
        ? err.message
        : 'Failed to fetch logins');
    }
    finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this login?'))
      return;

    try {
      await loginApi.deleteLogin(id);
      setLogins(logins().filter(login => login.id !== id));
    }
    catch (err) {
      setError(err instanceof Error
        ? err.message
        : 'Failed to delete login');
    }
  };

  onMount(() => {
    fetchLogins();
  });

  return (
    <div
      class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8"
    >
      <div
        class="sm:flex sm:items-center"
      >
        <div
          class="sm:flex-auto"
        >
          <h1
            class="text-2xl font-semibold text-gray-12"
          >
            User Logins
          </h1>
          <p
            class="mt-2 text-sm text-gray-9"
          >
            A list of all user logins in the system including their username, status, and password change information.
          </p>
        </div>
        <div
          class="mt-4 sm:mt-0 sm:ml-16 sm:flex-none"
        >
          <button
            class="inline-flex items-center justify-center rounded-lg border border-transparent bg-blue-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 sm:w-auto"
            type="button"
            onClick={() => navigate('/logins/create')}
          >
            Add Login
          </button>
        </div>
      </div>

      <Show
        when={error()}
      >
        <div
          class="mt-4 bg-red-50 p-4 rounded-lg"
        >
          <div
            class="flex"
          >
            <div
              class="flex-shrink-0"
            >
              <svg
                class="h-5 w-5 text-red-400"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path
                  clip-rule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                  fill-rule="evenodd"
                />
              </svg>
            </div>
            <div
              class="ml-3"
            >
              <h3
                class="text-sm font-medium text-red-800"
              >
                {error()}
              </h3>
            </div>
          </div>
        </div>
      </Show>

      <div
        class="mt-8 flex flex-col"
      >
        <div
          class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8"
        >
          <div
            class="inline-block min-w-full py-2 align-middle md:px-6 lg:px-8"
          >
            <div
              class="overflow-hidden shadow ring-1 ring-black ring-opacity-5 rounded-lg"
            >
              <table
                class="min-w-full divide-y divide-gray-6"
              >
                <thead>
                  <tr>
                    <th
                      class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-11 sm:pl-6"
                      scope="col"
                    >
                      Username
                    </th>
                    {/* 2FA Enabled column removed */}
                    <th
                      class="px-3 py-3.5 text-left text-sm font-semibold text-gray-11"
                      scope="col"
                    >
                      Password Last Changed
                    </th>
                    <th
                      class="px-3 py-3.5 text-left text-sm font-semibold text-gray-11"
                      scope="col"
                    >
                      Status
                    </th>
                    <th
                      class="relative py-3.5 pl-3 pr-4 sm:pr-6"
                      scope="col"
                    >
                      <span
                        class="sr-only"
                      >
                        Actions
                      </span>
                    </th>
                  </tr>
                </thead>
                <tbody
                  class="divide-y divide-gray-6"
                >
                  <Show
                    when={!loading()}
                    fallback={(
                      <tr>
                        <td
                          class="text-center py-4"
                          colspan="4"
                        >
                          Loading...
                        </td>
                      </tr>
                    )}
                  >
                    {logins().length === 0 ? (
                      <tr>
                        <td
                          class="text-center py-4 text-sm text-gray-9"
                          colspan="4"
                        >
                          No logins found
                        </td>
                      </tr>
                    ) : (
                      <For
                        each={logins()}
                      >
                        {
                          login => (
                            <tr
                              data-key={login.id}
                            >
                              <td
                                class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-11 sm:pl-6"
                              >
                                <a
                                  class="text-blue-600 hover:text-blue-800 hover:underline"
                                  href={`/logins/${login.id}/detail`}
                                  onClick={(e) => {
                                    e.preventDefault();
                                    navigate(`/logins/${login.id}/detail`);
                                  }}
                                >
                                  {login.username}
                                </a>
                              </td>
                              {/* 2FA Enabled column cell removed */}
                              <td
                                class="whitespace-nowrap px-3 py-4 text-sm text-gray-11"
                              >
                                {login.password_last_changed ? new Date(login.password_last_changed).toLocaleDateString() : '-'}
                              </td>
                              <td
                                class="whitespace-nowrap px-3 py-4 text-sm text-gray-11"
                              >
                                {login.status || 'Active'}
                              </td>
                              <td
                                class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-6"
                              >
                                <button
                                  class="text-blue-600 hover:text-blue-900 mr-4"
                                  onClick={() => navigate(`/logins/${login.id}/edit`)}
                                >
                                  Edit
                                </button>
                                <button
                                  class="text-red-600 hover:text-red-900"
                                  onClick={() => handleDelete(login.id!)}
                                >
                                  Delete
                                </button>
                              </td>
                            </tr>
                          )
                        }
                      </For>
                    )}
                  </Show>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Logins;

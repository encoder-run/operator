import { createBrowserRouter } from "react-router-dom";
import Sidebar from "../pages/nav";
import ErrorPage from "../pages/error";
import RepositoriesPage from "../pages/repositories";
import ModelsPage from "../pages/models";
import StoragePage from "../pages/storage";
import ModelDetailsPage from "../pages/model-details";
import CodeSearchPage from "../pages/search";
import RepositoryDetailsPage from "../pages/repository-details";
import StorageDetailsPage from "../pages/storage-details";
import PipelinesPage from "../pages/pipelines";
import PipelineDetailsPage from "../pages/pipeline-details";

export const router = createBrowserRouter([
    {
        path: "/",
        element: <Sidebar />,
        errorElement: <ErrorPage />,
        children: [
            {
                path: "search",
                element: <CodeSearchPage />,
            },
            {
                path: "api",
                element: <></>,
            },
            {
                path: "pipelines",
                element: <PipelinesPage />,
            },
            {
                path: "pipelines/:pipelineId",
                element: <PipelineDetailsPage />,
            },
            {
                path: "models",
                element: <ModelsPage />,
            },
            {
                path: "models/:modelId",
                element: <ModelDetailsPage />,
            },
            {
                path: "repositories",
                element: <RepositoriesPage />
            },
            {
                path: "repositories/:repositoryId",
                element: <RepositoryDetailsPage />
            },
            {
                path: "storage",
                element: <StoragePage />
            },
            {
                path: "storage/:storageId",
                element: <StorageDetailsPage />
            }
        ]
    },
]);
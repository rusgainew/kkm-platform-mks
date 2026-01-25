/**
 * Hook for loading companies with error handling
 */

import { useEffect, useState } from "react";
import { listCompanies, Company } from "@/lib/api/companies";

export const useCompanies = () => {
  const [companies, setCompanies] = useState<Company[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadCompanies = async () => {
      try {
        setIsLoading(true);
        setError(null);

        const response = await listCompanies(1, 100);

        console.log("[useCompanies] Raw response:", response);

        let companiesData: Company[] = [];

        // Check if response has data property
        if (response && typeof response === "object" && "data" in response) {
          companiesData = response.data;
          console.log(
            "[useCompanies] Extracted from data property:",
            companiesData
          );
        } else if (Array.isArray(response)) {
          companiesData = response;
          console.log("[useCompanies] Used as array:", companiesData);
        } else {
          console.log("[useCompanies] Unknown response format");
        }

        console.log("[useCompanies] Loaded companies:", companiesData);
        setCompanies(companiesData);
      } catch (err) {
        const errorMessage =
          err instanceof Error ? err.message : "Ошибка загрузки компаний";
        console.error("[useCompanies] Ошибка:", errorMessage);
        setError(errorMessage);
      } finally {
        setIsLoading(false);
      }
    };

    loadCompanies();
  }, []);

  return { companies, isLoading, error };
};
